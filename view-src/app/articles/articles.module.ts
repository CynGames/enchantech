import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ArticlesService} from "./articles.service";

@NgModule({
  declarations: [],
  imports: [
    CommonModule
  ],
  providers: [
    ArticlesService
  ],
})
export class ArticlesModule {
}
