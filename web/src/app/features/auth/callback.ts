import { ChangeDetectionStrategy, Component, inject, OnInit } from "@angular/core";
import {ActivatedRoute, Router} from "@angular/router";

import { AuthService } from "../../core/services/auth";

@Component({
  selector: "app-auth-callback",
  template: "<p>Autenticando...</p>",
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AuthCallback implements OnInit {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);
  private readonly route = inject(ActivatedRoute);

  ngOnInit(): void {
    const token = this.route.snapshot.queryParamMap.get('token');

    if (token){
      this.authService.setToken(token);
    }

    this.router.navigate(['/']);
  }

}
