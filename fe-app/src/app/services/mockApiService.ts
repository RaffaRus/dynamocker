import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { initialMockApiJson } from "@models/editor.model";
import { IMockApi, IModifiedMockApi } from "@models/mockApi.model";
import { BehaviorSubject, Observable } from "rxjs";

@Injectable()
export class MockApiService {
    
    private MOCK_API_SERVE_URL_BASE = 'http://localhost:8150/dynamocker/api'; 
    private MOCK_APIS = '/mock-apis'; 
    private MOCK_API = '/mock-api';
    
    private newMockApiSub : BehaviorSubject<IMockApi> = new BehaviorSubject<IMockApi>(initialMockApiJson)
    private refreshListSub : BehaviorSubject<null> = new BehaviorSubject<null>(null)

    constructor(private httpClient: HttpClient) { }

    newMockApiSelectedObservable() : Observable<IMockApi> {
        return this.newMockApiSub.asObservable()
    }

    selectMockApi(mockApi : IMockApi) {
        this.newMockApiSub.next(mockApi)
        console.log("item selected")
    }

    refreshListObservable() {
        return this.refreshListSub.asObservable()
    }
    
    refreshList() {
        console.log("someone requested a refresh of the list")
        this.refreshListSub.next(null)
    }

    getAllMockApis() : Observable<IMockApi[] > {
        let url = this.MOCK_API_SERVE_URL_BASE + this.MOCK_APIS
        return this.httpClient.get<IMockApi[] >(url)
    }

    deleteMockApi(mockApiName : string) : Observable<null> {
        let url = this.MOCK_API_SERVE_URL_BASE + this.MOCK_API + "/" + mockApiName
        return this.httpClient.delete<null>(url)
    }

    getMockApi(mockApiName : string) : Observable<IMockApi > {
        let url = this.MOCK_API_SERVE_URL_BASE + this.MOCK_API + "/" +  mockApiName
        return this.httpClient.get<IMockApi >(url)
    }
    
    postMockApi(mockApi : IMockApi) : Observable<null> {
        let url = this.MOCK_API_SERVE_URL_BASE + this.MOCK_API + "/" +  mockApi.name
        return this.httpClient.post<null>(url, JSON.stringify(mockApi))
    }
    
    patchMockApi(mockApi : IMockApi) : Observable<null> {
        let url = this.MOCK_API_SERVE_URL_BASE + this.MOCK_API + "/" +  mockApi.name
        return this.httpClient.patch<null>(url, JSON.stringify(mockApi))
    }

}
