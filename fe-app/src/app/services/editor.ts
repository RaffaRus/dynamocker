import { Injectable } from "@angular/core";

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
