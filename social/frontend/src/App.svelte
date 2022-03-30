<script>
  import Auth from './Auth.svelte'
  import Profile from './Profile.svelte'

  let state = 'auth';
</script>

{#if state === 'auth'}
  <Auth/>
{:else if state === 'submit'}
  <Profile/>
{:else if state === 'list' }
  <div></div>
{/if}
