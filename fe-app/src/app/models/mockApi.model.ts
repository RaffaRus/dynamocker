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

// TODO: allow also objects to used under each element of the IResponse
// The editor of the ui at the moment returns a warning because an object
// is expected under each elements of the IResponse
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