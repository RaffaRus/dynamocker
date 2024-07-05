export interface  IMockApi {
    Name :string
    URL: string
    Added: Date
    LastModified: Date
    Responses: IResponse
}

interface IResponse {
    Get?: JSON
    Patch?: JSON
    Delete?: JSON
    Post?: JSON
}

export const dummyMockApi : IMockApi = {
    Added: new Date(),
    LastModified: new Date(),
    Name: "",
    Responses : {
      Delete: JSON.parse("{\"delete\":\"delete\"}"),
      Get: JSON.parse("{\"get\":\"get\"}"),
      Patch: JSON.parse("{\"patch\":\"patch\"}"),
      Post: JSON.parse("{\"post\":\"post\"}"),
    },
    URL: ""
  }
