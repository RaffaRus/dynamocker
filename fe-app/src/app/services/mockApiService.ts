import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { MockApi, dummyMockApi } from "@models/mockApi";
import { BehaviorSubject, Observable } from "rxjs";

@Injectable()
export class MockApiService {
    
    private MOCK_API_SERVE_URL_BASE = 'http://localhost:<PORT>/dynamocker/api'; 
    private MOCK_APIS = '/mock-apis'; 
    private MOCK_API = '/mock-apis';
    
    private newMockApiSub : BehaviorSubject<MockApi> = new BehaviorSubject<MockApi>(dummyMockApi)
    private refreshListSub : BehaviorSubject<null> = new BehaviorSubject<null>(null)

    constructor(private httpClient: HttpClient) { }

    newMockApiSelectedObservable() : Observable<MockApi> {
        return this.newMockApiSub.asObservable()
    }

    selectMockApi(mockApi : MockApi) {
        this.newMockApiSub.next(mockApi)
    }

    refreshListObservable() {
        return this.refreshListSub.asObservable()
    }
    
    refreshList() {
        this.refreshListSub.next(null)
    }

    getAllMockApis() : Observable<MockApi[] > {
        let url = this.MOCK_API_SERVE_URL_BASE + this.MOCK_APIS
        return this.httpClient.get<MockApi[] >(url)
    }

    deleteAllMockApis(mockApiName : string) : Observable<null > {
        let url = this.MOCK_API_SERVE_URL_BASE + this.MOCK_API + mockApiName
        return this.httpClient.delete<null >(url)
    }

    getMockApi(mockApiName : string) : Observable<MockApi > {
        let url = this.MOCK_API_SERVE_URL_BASE + this.MOCK_API + mockApiName
        return this.httpClient.get<MockApi >(url)
    }
    
    postMockApi(mockApi : MockApi) : Observable<null > {
        let url = this.MOCK_API_SERVE_URL_BASE + this.MOCK_API
        return this.httpClient.post<null >(url, JSON.stringify(mockApi))
    }
    
    patchMockApi(mockApi : MockApi) : Observable<null > {
        let url = this.MOCK_API_SERVE_URL_BASE + this.MOCK_API
        return this.httpClient.patch<null >(url, JSON.stringify(mockApi))
    }

}
