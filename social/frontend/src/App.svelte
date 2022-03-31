<script>
  import { state } from './stores.js'
  import Auth from './Auth.svelte'
  import Profile from './Profile.svelte'
  import Profiles from './Profiles.svelte'

  let pageState;

  state.subscribe(value => {
    pageState = value;
  });
</script>

{#if pageState === 'auth'}
  <Auth/>
{:else if pageState === 'submit'}
  <Profile/>
{:else if pageState === 'profiles' }
  <Profiles/>
{/if}
