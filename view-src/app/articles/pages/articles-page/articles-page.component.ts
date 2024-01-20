import {Component, OnInit} from '@angular/core';
import {ArticlesService} from "../../articles.service";

@Component({
  selector: 'app-articles-page',
  standalone: true,
  imports: [],
  templateUrl: './articles-page.component.html'
})
export class ArticlesPageComponent implements OnInit{
  constructor(
    private readonly articlesService: ArticlesService,
  ) {
  }

  ngOnInit() {
    this.articlesService.getTest().subscribe({
      next: (data) => {
        console.log('data', data);
      },
      error: (err) => {
        console.log('err', err);
      }
    });
  }
}
