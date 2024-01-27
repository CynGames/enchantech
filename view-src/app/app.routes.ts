import { Routes } from '@angular/router';
import {ROGUINNIEComponent} from "../rogu-innie/rogu-innie.component";
import {ArticlesPageComponent} from "../articles/articles-page.component";

export const routes: Routes = [
  {
    component: ROGUINNIEComponent,
    path: 'login'
  },
  {
    component: ArticlesPageComponent,
    path: ''
  }
];
