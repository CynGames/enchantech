import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";

@Injectable()
export class ArticlesService {
  constructor(private http: HttpClient) { }


  getTest() {
    return this.http.get('/api/test');
  }

}
