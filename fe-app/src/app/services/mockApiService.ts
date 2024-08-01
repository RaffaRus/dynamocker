import { HttpClient, HttpErrorResponse } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { IMockApi, ResourceObject } from "@models/mockApi.model";
import { BehaviorSubject, map, Observable, tap, throwError } from "rxjs";

@Injectable()
export class MockApiService {

    private MOCK_API_SERVE_URL_BASE = 'http://localhost:8150/dynamocker/api';
    private MOCK_APIS = '/mock-apis';
    private MOCK_API = '/mock-api';

    private newMockApiSub: BehaviorSubject<ResourceObject> = new BehaviorSubject<ResourceObject>(new ResourceObject())
    private refreshListSub: BehaviorSubject<null> = new BehaviorSubject<null>(null)

    constructor(private httpClient: HttpClient) { }

    newMockApiSelectedObservable(): Observable<ResourceObject> {
        return this.newMockApiSub.asObservable()
    }

    selectMockApi(resObj: ResourceObject) {
        this.newMockApiSub.next(resObj)
    }

    refreshListObservable() {
        return this.refreshListSub.asObservable()
    }

    refreshList() {
        this.refreshListSub.next(null)
    }

    getAllMockApis(): Observable<ResourceObject[]> {
        let url = this.MOCK_API_SERVE_URL_BASE + this.MOCK_APIS
        return this.httpClient.get<ResourceObject[]>(url)
    }

    deleteMockApi(mockApiUuid: number): Observable<null> {
        let url = this.MOCK_API_SERVE_URL_BASE + this.MOCK_API + "/" + mockApiUuid
        return this.httpClient.delete<null>(url)
    }

    getMockApi(mockApiUuid: string): Observable<ResourceObject> {
        let url = this.MOCK_API_SERVE_URL_BASE + this.MOCK_API + "/" + mockApiUuid
        return this.httpClient.get<ResourceObject>(url)
    }

    postMockApi(mockApi : IMockApi) : Observable<null> {
        let url = this.MOCK_API_SERVE_URL_BASE + this.MOCK_API
        return this.httpClient.post<null>(url, JSON.stringify(mockApi))
    }

    putMockApi(resObj: ResourceObject): Observable<null> {
        var resList: ResourceObject[] = [];
        this.getAllMockApis().subscribe({
            next: (value) => {
                console.log(value)
                resList = value
            },
            error: (err) => {
                console.error(err)
            }
        })
        for (let i = 0; i < resList.length; i++) {
            const element = resList[i];
            if (element.data.name == resObj.data.name) {
                throw new HttpErrorResponse({
                    error: 'found another mockApi using the same name',
                    status: 500,
                    statusText: 'Warning'
                })
            } else if (element.data.url == resObj.data.url) {
                throw new HttpErrorResponse({
                    error: 'found another mockApi using the same url',
                    status: 500,
                    statusText: 'Warning'
                });
            }
        }
        let url = this.MOCK_API_SERVE_URL_BASE + this.MOCK_API + "/" + resObj.id
        return this.httpClient.put<null>(url, JSON.stringify(resObj.data))
    }

}
