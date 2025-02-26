(function () {
  console.log("Incoming headers are", JSON.stringify(request.headers))

  return {
    status: 200,
    body: {
      users: {
        list: [
          {
            id: "1",
            name: "John",
            surname: "Doe"
          }
        ]
      }
    }
  }
})()
