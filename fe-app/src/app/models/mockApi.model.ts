import { initialMockApiJsonString } from "./editor.model"

export interface  IResourceObject {
    id : number
    type : string
    data: IMockApi
}

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

export class ResourceObject implements IResourceObject {
    id = -1
    type = "mockApi"
    data = JSON.parse(initialMockApiJsonString)
  }