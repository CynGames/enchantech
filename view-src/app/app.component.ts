import {Component} from '@angular/core';
import {RouterOutlet} from '@angular/router';
import {AuthModule} from "../auth/auth.module";
import {NavbarComponent} from "../navbar/navbar.component";

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, AuthModule, NavbarComponent],
  templateUrl: './app.component.html'
})
export class AppComponent {
}
