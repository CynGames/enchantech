import { Injectable } from '@angular/core';
import {HttpClient, HttpParams} from "@angular/common/http";
import {Observable} from "rxjs";
import {Article} from "./article";

@Injectable()
export class ArticlesService {
  constructor(private http: HttpClient) { }


  fetchArticles(filter: FetchArticlesFilter) {

    const params  = new HttpParams({
      fromObject: Object.entries(filter).reduce((acc, [key, value]) => {
        acc[key] = String(value);
        return acc;
      }, {} as { [key: string]: string })
    })

    return this.http.get('/api/articles', {
      params,
    }) as Observable<Article[]>;
  }
}


export interface FetchArticlesFilter {
  limit: number,
  offset: number,
}
