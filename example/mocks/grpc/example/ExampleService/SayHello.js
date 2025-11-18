(function () {
  console.log("Incoming metadata are", JSON.stringify(request.metadata))

  let err = request.metadata['test-case-error']
  if (err) {
    return {
      error: {
        code: 3,
        message: "Invalid argument"
      }
    }
  }

  return {
    body: {
      message: `Hello, ${request.body.name}, your role is ${request.body.role}`,
      details: {
        code: 1,
        status: "OK"
      }
    }
  }
})()
