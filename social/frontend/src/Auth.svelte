<script>
  import { state } from './stores.js';

  let nickname = "";
  let password = "";

  function login(nickname, password) {
    auth('login', nickname, password);
    state.update(_ => "profiles");
  }

  function signup(nickname, password) {
    auth('signup', nickname, password);
    state.update(_ => "submit");
  }

  function auth(path, nickname, password) {
    fetch('/api/v1/' + path, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({nickname: nickname, password: password})
    }).then(async response => {
      const object = await response.json();
      sessionStorage.setItem("access_token", object.access_token);
    });
  }
</script>

<section class="hero is-success is-fullheight">
  <div class="hero-body">
    <div class="container has-text-centered">
      <div class="column is-4 is-offset-4">
        <div class="box">
          <form>
            <div class="field">
              <div class="control">
                <input bind:value={nickname}
                  class="input is-large"
                  type="nickname"
                  placeholder="Your nickname"
                />
              </div>
            </div>

            <div class="field">
              <div class="control">
                <input bind:value={password}
                  class="input is-large"
                  type="password"
                  placeholder="Your Password"
                />
              </div>
            </div>
            <button type="button" class="button is-block is-info is-large is-fullwidth mb-2" on:click={login(nickname, password)}>
              Login <i class="fa fa-sign-in" aria-hidden="true"></i>
            </button>
            <button type="button" class="button is-block is-info is-large is-fullwidth mb-2" on:click={signup(nickname, password)}>
              Signup <i class="fa fa-sign-in" aria-hidden="true"></i>
            </button>
          </form>
        </div>
      </div>
    </div>
  </div>
</section>

<style>
.hero.is-success {
  background: #F2F6FA;
}
.box {
  margin-top: 5rem;
}
input {
  font-weight: 300;
}

.field{
  padding-bottom: 10px;
}

.fa{
  margin-left: 5px;
}
</style>
