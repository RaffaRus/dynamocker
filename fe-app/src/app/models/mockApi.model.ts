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

export interface  IModifiedMockApi {
    Name :string
    URL: string
    Added: Date
    LastModified: Date
    Responses: IModifiedResponse
}

interface IModifiedResponse {
    Get?: string
    Patch?: string
    Delete?: string
    Post?: string
}