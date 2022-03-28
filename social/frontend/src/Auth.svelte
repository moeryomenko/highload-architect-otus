<script>
  export let operation;
  export let titleOperation = "";

  let nickname = "";
  let password = "";

  $: title = capitalizeFirstLetter(titleOperation)

  function capitalizeFirstLetter(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
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
            <button type="button" class="button is-block is-info is-large is-fullwidth" on:click={operation(nickname, password)}>
              {title} <i class="fa fa-sign-in" aria-hidden="true"></i>
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
p {
  font-weight: 700;
}
p.subtitle {
  padding-top: 1rem;
}

.login-hr{
  border-bottom: 1px solid black;
}

.has-text-black{
  color: black;
}

.field{
  padding-bottom: 10px;
}

.fa{
  margin-left: 5px;
}
</style>
