export interface  IMockApi {
    name :string
    url: string
    added: Date
    responses: IResponse
}

interface IResponse {
    get?: JSON
    patch?: JSON
    delete?: JSON
    post?: JSON
}

export interface  IModifiedMockApi {
    name :string
    url: string
    added: Date
    responses: IModifiedResponse
}

interface IModifiedResponse {
    get?: string
    patch?: string
    delete?: string
    post?: string
}