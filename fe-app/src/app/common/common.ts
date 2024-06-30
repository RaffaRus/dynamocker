export function isDefined<Type>(elem : Type) : boolean {
    if (elem !== null && elem !== undefined) {
       return true
    } else {
       return false 
    }
}