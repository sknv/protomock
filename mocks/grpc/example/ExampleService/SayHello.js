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
