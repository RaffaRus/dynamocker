import { HttpHeaders } from "@angular/common/http"

export function isDefined<Type>(elem : Type) : boolean {
    if (elem !== null && elem !== undefined) {
       return true
    } else {
       return false 
    }
}

export const corsHeaders : HttpHeaders = new HttpHeaders({
   'Content-Type': 'application/json',
   'Accept': 'application/json',
   'Access-Control-Allow-Headers': 'Content-Type',
 })