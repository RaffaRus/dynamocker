import { Time } from "@angular/common"

export interface  MockApi {
    Name :string
    URL: string
    Added: Time
    LastModified: Time
    Responses: Response
}

interface Response {
    Get: JSON
    Patch: JSON
    Delete: JSON
    Post: JSON
}

