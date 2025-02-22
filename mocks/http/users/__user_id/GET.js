function main() {
  let name = request.headers["Test-Case-Name"] ?? "John"

  return {
    status: 200,
    headers: null,
    body: {
      user: {
        id: `${request.params.user_id}`,
        name: name,
        surname: "Doe"
      }
    }
  }
}

main()
