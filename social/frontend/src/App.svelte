<script>
  import Auth from './Auth.svelte'

  let token = localStorage.getItem("access_token") || "";

  function login(nickname, password) {
    // localStorage.setItem("access_token", token)
  }

  function signup(nickname, password) {
    // localStorage.setItem("access_token", token)
  }

  function activeTab(tab) {
  }
</script>

<Auth titleOperation="sign up" operation={signup}/>

<style>
#tabs-with-content .tabs:not(:last-child) {
  margin-bottom: 0;
}

#tabs-with-content .tab-content {
  padding: 1rem;
  display: none;
}

#tabs-with-content .tab-content.is-active {
  display: block;
}
</style>
