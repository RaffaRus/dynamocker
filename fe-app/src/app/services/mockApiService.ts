import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { ResourceObject } from '@models/mockApi.model';
import { IMockApi, IResourceObject } from "@models/mockApi.model";
import { BehaviorSubject, Observable } from "rxjs";

@Injectable()
export class MockApiService {
    
    private MOCK_API_SERVE_URL_BASE = 'http://localhost:8150/dynamocker/api'; 
    private MOCK_APIS = '/mock-apis'; 
    private MOCK_API = '/mock-api';
    
    private newMockApiSub : BehaviorSubject<IResourceObject> = new BehaviorSubject<IResourceObject>(new ResourceObject())
    private refreshListSub : BehaviorSubject<null> = new BehaviorSubject<null>(null)

    constructor(private httpClient: HttpClient) { }

    newMockApiSelectedObservable() : Observable<IResourceObject> {
        return this.newMockApiSub.asObservable()
    }

    selectMockApi(resObj: IResourceObject) {
        this.newMockApiSub.next(resObj)
    }

    refreshListObservable() {
        return this.refreshListSub.asObservable()
    }
    
    refreshList() {
        this.refreshListSub.next(null)
    }

    getAllMockApis() : Observable<IResourceObject[] > {
        let url = this.MOCK_API_SERVE_URL_BASE + this.MOCK_APIS
        return this.httpClient.get<IResourceObject[] >(url)
    }

    deleteMockApi(mockApiUuid :number) : Observable<null> {
        let url = this.MOCK_API_SERVE_URL_BASE + this.MOCK_API + "/" + mockApiUuid
        return this.httpClient.delete<null>(url)
    }

    getMockApi(mockApiUuid : string) : Observable<IResourceObject > {
        let url = this.MOCK_API_SERVE_URL_BASE + this.MOCK_API + "/" +  mockApiUuid
        return this.httpClient.get<IResourceObject>(url)
    }
    
    postMockApi(mockApi : IMockApi) : Observable<null> {
        let url = this.MOCK_API_SERVE_URL_BASE + this.MOCK_API
        return this.httpClient.post<null>(url, JSON.stringify(mockApi))
    }
    
    putMockApi(resObj : IResourceObject) : Observable<null> {
        let url = this.MOCK_API_SERVE_URL_BASE + this.MOCK_API + "/" +  resObj.id
        return this.httpClient.put<null>(url, JSON.stringify(resObj.data))
    }

}
