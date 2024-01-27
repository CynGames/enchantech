import {Component, OnInit} from '@angular/core';
import {ArticlesService} from "./articles.service";
import {Article} from "./article";
import {CommonModule} from "@angular/common";

@Component({
  selector: 'app-articles-page',
  standalone: true,
  imports: [CommonModule],
  providers: [ArticlesService],
  templateUrl: './articles-page.component.html'
})
export class ArticlesPageComponent implements OnInit {
  constructor(
    private readonly articlesService: ArticlesService,
  ) {
  }

  protected articles: Article[] = [];

  ngOnInit() {
    this.fetchMoreArticles();
  }

  protected fetchMoreArticles() {
    // TODO maybe theres a way to take advantage of observables to achieve pagination
    console.log('fetchMoreArticles')
    this.articlesService.fetchArticles({
      limit: 3,
      offset: this.articles.length,
    }).subscribe({
      next: (data) => {
        this.articles = this.articles.concat(data);
      },
      error: (err) => {
        console.log('err', err);
      }
    });
  }

}

