import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { IMockApi, dummyMockApi } from "@models/mockApi";
import { BehaviorSubject, Observable } from "rxjs";

@Injectable()
export class EditorService {

    private _unsavedModifications : boolean = false

    constructor(
        
    ) { }

    public get unsavedModifications() {
        return this._unsavedModifications
    }

    public set unsavedModifications(bool : boolean) {
        this._unsavedModifications = bool
    }

}
