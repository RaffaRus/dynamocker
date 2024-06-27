export function isDefined<Type>(elem : Type) : boolean {
    if (elem !== null && elem !== undefined) {
       return false
    } else {
       return true 
    }
}