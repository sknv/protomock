# protomock

A simple mock server that supports HTTP and gRPC protocols. The main idea is to define mocks using JavaScript making them highly customizable in integration tests.

## What is in the box

protomock can serve two types of mocks:

- HTTP requests receiving and responding JSON
- gRPC unary calls

## To Be Done

- Add gRPC streaming support
- Add gRPC unions support

## Configuration

The only required configuration is a `.yaml` file with self-explanatory sections. You can provide a path to a configuration file via `-c` flag or omit one and use the default path `./configs/protomock.yaml`.

## Mock definition

protomock follows the "convention over configuration" approach to define mocks. That means you only have to place your mock files in specific folders and protomock will do the rest.

### HTTP mock definition

HTTP mocks are constructed based on a folder tree.

#### GET request to /hello-world

- Create a directory `${HTTP_SERVER_MOCKSDIR}/hello-world`
- Create a `GET.js` file inside it with your JS code

#### POST request to /users

- Create a directory `${HTTP_SERVER_MOCKSDIR}/users`
- Create a `POST.js` file inside it with your JS code

#### GET request to /users/:user_id

- Create a directory `${HTTP_SERVER_MOCKSDIR}/users/__user_id` (double underscore prefix)
- Create a `GET.js` file inside it with your JS code

Similarly you can create `PUT.js`, `DELETE.js` etc in your path. For wildcard params use directory name as `__param` (double underscore prefix).

#### HTTP request and response

Inside a mock you have access to the following request parameters:

- URL parameters
- Headers
- JSON body

```js
let request = {
  params: { // URL parameters object
    ...
  },
  headers: { // Headers object
    ...
  },
  body: { // JSON body
    ...
  }
}
```

`request` object is implicitly injected to your script.

The response has the following structure:
```js
let response = {
  status:  200, // HTTP status code
  body: { // JSON body
    ...
  }
}
```

You also can log any information via `console.log` function.

A sample JS mock file is presented below:

```js
(function () {
  console.log("Incoming headers are", JSON.stringify(request.headers))

  let name = request.headers["Test-Case-Name"] ?? "John"

  return {
    status: 200,
    body: {
      user: {
        id: `${request.params.user_id}`,
        name: name,
        surname: "Doe"
      }
    }
  }
})()
```

### gRPC mock definition

To define gRPC mocks you need to provide both `.proto` definition and `.js` mock files. For example, lets define an `ExampleService` inside an `example` proto package.

- Create a directory `${GRPC_SERVER_MOCKSDIR}/example`, where `example` is a package name and all the files belong to the package should be placed inside
- Place a `.proto` definition inside, for example:

```proto
syntax = "proto3";

package example;

service ExampleService {
  rpc SayHello (HelloRequest) returns (HelloResponse);
}

message HelloRequest {
  string name = 1;
}

message HelloResponse {
  string message = 1;
  Details details = 2;
}

message Details {
  int32 code = 1;
  string status = 2;
}
```

- Create a folder tree in form of `ServiceName/MethodName.js`, e.g. `ExampleService/SayHello.js`

Now all the gRPC requests to `example.ExampleService.SayHello` method will use `SayHello.js` code to build a response.

#### gRPC request and response

Inside a mock you have access to the following request parameters:

- Metadata
- Proto body

```js
let request = {
  metadata: { // Metadata object
    ...
  },
  body: { // Proto body
    ...
  }
}
```

`request` object is implicitly injected to your script.

The response has the following structure:
```js
let response = {
  body: { // Proto body
    ...
  }
}
```

You also can log any information via `console.log` function.

A sample JS mock file is presented below:

```js
(function () {
  console.log("Incoming metadata are", JSON.stringify(request.metadata))

  return {
    body: {
      message: `Hello, ${request.body.name}`,
      details: {
        code: 1,
        status: "OK"
      }
    }
  }
})()
```
