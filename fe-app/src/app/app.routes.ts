import { Routes } from '@angular/router';
import { BackgroundComponent } from './components/background/background.component';

export const routes: Routes = [
    {
        path: '/',
        component: BackgroundComponent,
    },
    {
        path: '**',
        redirectTo: '/',
    },
]