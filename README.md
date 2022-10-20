# README

This is the tech task I was asked to do as part of my selection process on ESL FACEIT Group. It’s a basic microservice
handling a CRUD on a user entity, complying with the specified requirements.

# How to run it

I tried to make running the microservice as easy as possible, in order to avoid problems on its review. All is needed to
run this project is having Docker installed and available on the PATH. Note that all the containers are expected to run
in linux/amd64 images.

In the Makefile we can find a bunch of commands. In order to review this project all is needed is three of them:

- `run.compose` This deploys a Docker Compose template that builds and runs our microservice along with a local MongoDB
  instance.
- `test.unit` This runs unit testing on the `internal` directory (as it’s the only directory where unit testing makes
  sense)
- `test.e2e` This deploys a Docker Compose template including our microservice, a MongoDB instance, and a custom
  container executing the E2E tests located in `test/e2e`

# Tech stack

As required, the microservice is written in Golang. I chose not to use any framework or similar libraries. Considering
the options, I preferred to keep the dependency tree slim. Also, I didn’t think that I really needed one, as I could
achieve the task requirements with only a few packages.

The only libraries I used are, on one hand, the required drivers, like MongoDB and gRPC, and on the other, small helper
utilities, like logrus for logging, testify for the unit and e2e tests, a library with helper methods for errors, and
so.

Apart from the microservice dependencies I used the protoc compiler and its gRPC plugin to generate code from the
protobuf definitions, and all the processes, like compiling, testing, or code generation run on Docker containers. I
made a special effort to comply with this, given that the people who will review this project will be alien to its
implementation. So this way I can be sure that what works in my computer, does in theirs.

# Repository structure

When it comes to Golang microservices, I usually stick to the same layout for the repository, based on
this [non-official standard](https://github.com/golang-standards/project-layout).

Summing things up on how it works, it's a monorepo layout, great for situations where you have several microservices. In
this case we only have one, but it could be easily escalated. The code distributes on these directories:

- `internal/app/<service>` Most of the code for our microservices is located here. Here, each service can have its
  own `domain` layer (for multiple domain entities if needed), required `ports` exposing one or multiple APIs (gRPC,
  REST, GraphQL…), and any needed `adapters`, to connect to a database or to other (micro)services. We can also find
  the `app` directory, holding our CQRS handlers and logic, and the `service` directory, where we assemble the
  dependency between the former and the adapters.
- `internal/pkg` This is where we store anything that might be useful and reusable for any microservice in the monorepo.
  For instance, I kept some MongoDB helper types there, and the code to generate a new gRPC server with some common
  configuration.
- `cmd/<service>` Here we hold the entrypoints for our microservices, keeping them as simple and slim as possible.
- `pkg` Here we keep anything that should be available to import for other external applications, in this case, the gRPC
  autogenerated stubs.

There are some other directories related to testing, deployment, and scripting:

- `api` We keep our Protocol Buffers definition files here. If our microservices had any other kind of API definition
  files, like OpenAPI, we would keep them here too.
- `build` The Dockerfiles needed to build or test our microservices are located here.
- `deploy` Any deployment templates or configuration should be kept here. In this case, we only have a couple of Docker
  Compose templates, but we could have other kinds, like k8s deployment templates.
- `examples` As the name indicates, examples on how to use our application.
- `scripts` Scripts for processes called from our Makefile, keeping it cleaner and slimmer.
- `test` Testing-related files. In this case, the files containing the E2E tests, and some autogenerated mocks.

# Service architecture

To develop this service I organised the code in a Clean Architecture, along with a CQRS pattern. I chose this because,
from what I’ve seen in the past, it produces a clean and highly scalable codebase.

Following the path a request to our microservice makes, we start with the **ports**. These process the incoming requests
to our service. We have as many ports as needed, each one handling one type of API. In this case, we have one only one
port, the one that handles our gRPC endpoints. The job of a port is to parse the request and call the required commands
and queries to produce the corresponding response. We try to keep the ports simple, and delegate most of the work to the
inner layers of our architecture.

A port can call one or multiple commands and queries to generate the desired behaviour. One example of an endpoint that
calls multiple of them -a command and then a query- is the UpdateUser endpoint, as the UpdateUser command doesn’t return
anything, and we usually want to return that data to the client on the same request. We could also avoid that, and
transfer that responsibility to the client, making the endpoints more simple, and even enabling us to make the requests
asynchronous (for instance, just returning an ack along with an id that would serve to trace when the requested process
has finished)

**Commands** and **queries** can be seen as the “orchestrators” of our processes. They receive the required data to
execute a specific action, interact with the domain, and call the adapters needed to complete their task. As expected,
queries are meant to retrieve data, while commands are meant to modify it. In fact, commands shouldn’t return any data
at all. The only exception I make for this is to return ids generated on the command handler, because it’s mostly needed
to access the data later.

Something that may not be really needed, but I like to do, is having specific struct types for the query responses, as
this gives them more flexibility on what the ports can receive back.

As stated, sometimes the queries and the commands interact with the **domain** layer, always using its exposed methods
in order to modify its state. Validations on whether a change is acceptable or not should be made in these methods, not
on the handlers.

At last, we have the **adapters**. These connect our microservices to external services in order to perform its own
requests when needed. Command and query handlers need to be provided with the adapters required to perform their duty on
their constructors, but we do this using a generic interface that the adapter can implement later on, this way we
decouple its implementation, and gain the ability to replace it easily if ever needed.

I like to classify my adapters in three types: **repositories**, for the ones that handle interaction with a database, *
*clients**, for the ones that make requests to other services, and **producers**, for the ones that publish events to a
message queue.

## Scalability

As I explained in the “Repository Structure” section, this monorepo can be easily scaled to hold several microservices
coexisting together, sharing reusable code and coordinating their deployments. you just need to add the required code
inside `internal/app`, create an entrypoint, and add the required scripts and build templates.

When it comes to scaling the ability of this kind of microservices to communicate with each other (or even with external
services), what I usually see as convenient is implementing an event-driven architecture, using a message queue (Apache
Kafka for example can be really powerful) to publish and consume events informing of the specific actions being
executed. This results in really clean interactions, as they would be completely asynchronous, while inverting the
dependency between the microservices: the responsibility of taking action in response to a specific event is no longer
placed on its producer, but on its consumer.

# Testing

I implemented two kinds of tests for this project, unit testing, and E2E testing.

The E2E tests are pretty simple, I try to simulate common use cases (create-read-update-remove), along with some
potential error cases.

The unit testing is a bit more complex, as some functions have been left out because of their specially specific
purpose. On the code sensible to test, I tried to get 100% of coverage. The overall coverage I got is 99.7%, with that
0.3% being a single line in `internal/pkg/utils/grpc_utils/mappers.go`, as testing it would have required an invalid
value in a Protocol Buffers message, which then would have required me either to “publish” a whole new version of the
proto definition, or craft a specific malicious protobuf body.

One thing I’d like to point out is the difficulty of mocking the official MongoDB Go driver, as they don’t use
interfaces to expose the database methods. This makes mocking them a bit difficult: I had to implement something kind of
like “wrappers” for the methods I needed, along with my own custom interfaces that would be the ones used in the
adapters and the ones I would mock for the tests. It feels a bit overengineered, but I really wanted to make an effort
to reach as much coverage as possible, so I decided to go with it. It can be seen in `internal/pkg/helper/mongo_helper`.

# Custom decisions I made

While performing this task, I had to make some decisions based on my own judgement, as they were not considered in the
task description.

## Hashing the passwords before storing them

This one was really easy for me to make. The description didn’t ask for the passwords to be protected in any form, but
from my experience, I know security needs to be a concern in all the stages of software development. I’d never store
user passwords in plain text, so I applied that same judgement here.

As stated in a code comment, I chose bcrypt as the hashing function knowing that there’s a more modern algorithm that is
starting to take over for this kind of use cases, Argon2id, in fact, is the top recommendation by the OWASP foundation
guidelines. I did this because I did a bit of research and there’s some debate on which algorithm is stronger to
brute-force attack (it seems bcrypt may be stronger when it comes to GPU-based attacks). So, knowing this, I decided to
stay with bcrypt, as it’s still strong, also recommended by OWASP, and has been battle-tested for a longer time.

## Not using any Go framework

As I stated previously, I chose to not use any specific Golang framework for this task. I only used some needed drivers
and a bunch of small libraries and utilities. There’s two main reasons I made this decision:

First of all, I wanted to keep the dependency tree slim. It can improve performance, it reduces the surface for supply
chain attacks (which are becoming more common), and you don’t get tied to specific frameworks, something that might
become a problem in the future as the functional and technical requirements grow.

Additionally, I tend to not use web frameworks in Go, because I don't usually feel like I need them. That's something I
appreciate about the language, the built-in packages can be really powerful. In JavaScript, for example, you end up
needing a specific library for a lot of things, so that wouldn't be a realistic approach, but in Golang it's viable.

## Fetching an entity before updating it

This is something specific to the UpdateUser command. The way it’s currently implemented, the handler receives the
command, fetches the user current state from the repository, executes the update (via the corresponding domain method),
and saves the full updated entity. At first glance, it may seem obvious that the same result could be achieved by just
sending an update query (`UPDATE` statement in SQL, `updateOne` in MongoDB…), and not only would it be more simple, but
also would reduce the latency.

However, the architecture we’re implementing requires all updates to go through our domain layer. This way, we can
ensure a proper validation of the changes, and avoid ending up with our domain in an unexpected state, specially as our
application grows. So, the best thing to do is to grab all the entity data, and implement methods that handle specific
situations, ensuring the required logic is applied.

## Health check

I implemented the health check following [the standard](https://github.com/grpc/grpc/blob/master/doc/health-checking.md)
defined by gRPC. This way, we can use standard probes to check our service is healthy and running.

However, the standard defines two methods for the health check server: a Check method, and a Watch method, and I left
the latter unimplemented. The thing is, the way the standard definition is intended to work is to have an in-memory
“registry” of the status of the different dependencies of the microservice, and it states the following about the Watch
method:

> A client can call the`Watch`method to perform a streaming health-check. The server will immediately send back a
> message indicating the current serving status. It will then subsequently send a new message whenever the service's
> serving status changes.

I didn’t want the health check to work like that, as it added a lot of complexity and involved the risk of the in-memory
registry not being accurate. I wanted the microservice to relay the check requests to the involved services in real
time. So I implemented the Check method only, following that logic, and did without the Watch method.
