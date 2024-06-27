export interface  MockApi {
    Name :string
    URL: string
    Added: Date
    LastModified: Date
    Responses: Response
}

interface Response {
    Get: JSON
    Patch: JSON
    Delete: JSON
    Post: JSON
}

export const dummyMockApi : MockApi = {
    Added: new Date(),
    LastModified: new Date(),
    Name: "",
    Responses : {
      Delete: JSON.parse(""),
      Get: JSON.parse(""),
      Patch: JSON.parse(""),
      Post: JSON.parse(""),
    },
    URL: ""
  }