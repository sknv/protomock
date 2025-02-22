function main() {
  let name = request.body.name ?? "John"
  let surname = request.body.surname ?? "Doe"

  return {
    status: 201,
    headers: null,
    body: {
      users: {
        list: [
          {
            id: "1",
            name: name,
            surname: surname
          }
        ]
      }
    }
  }
}

main()
