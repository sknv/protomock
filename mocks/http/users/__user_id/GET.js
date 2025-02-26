(function () {
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
