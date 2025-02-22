function main() {
  return {
    status: 200,
    headers: null,
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
}

main()
