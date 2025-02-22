function main() {
  return {
    status: 200,
    headers: null,
    body: {
      users: {
        list: [
          {
            name: "John",
            surname: "Doe"
          }
        ]
      }
    }
  }
}

main()
