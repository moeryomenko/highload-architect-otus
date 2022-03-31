<script>
  import { state } from './stores.js';

  let profile = {
    first_name: "",
    last_name:  "",
    age:        0,
    gender:     "",
    city:       "",
    interests:  [],
  };

  let genders = [
    {id: 1, text: "male"},
    {id: 2, text: "female"},
  ];

  let selected_gender;

  function mapToGender() {
    profile.gender = selected_gender.text;
  };

  function updateInterests() {
    profile.interests = [];
    for (let element of document.getElementsByClassName('tags-input')) {
      for (let item of element.children){
        if (item.classList.contains('tag')) {
          profile.interests.push(item.dataset.value);
        }
      }
    };
  }

  function submitProfile(profile) {
    updateInterests();
    mapToGender();

    fetch('/api/v1/profile', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${sessionStorage.getItem('access_token')}`,
       },
      body: JSON.stringify(profile),
    });

    state.update(_ => 'profiles');
  }
</script>

<div class="container is-center">
  <div class="field">
    <label class="label">First name</label>
    <div class="control">
      <input bind:value={profile.first_name} class="input" type="text">
    </div>
  </div>

  <div class="field">
    <label class="label">Last name</label>
    <div class="control">
      <input bind:value={profile.last_name} class="input" type="text">
    </div>
  </div>

  <div class="field">
    <label class="label">Age</label>
    <div class="control">
      <input bind:value={profile.age} class="input" type="number">
    </div>
  </div>

  <div class="field">
    <label class="label">Gender</label>
    <p class="control has-icons-left">
      <span class="select">
        <select bind:value={selected_gender}>
        {#each genders as gender}
          <option value={gender}>{gender.text}</option>
        {/each}
        </select>
      </span>
      <span class="icon is-small is-left">
        <i class="fas fa-globe"></i>
      </span>
    </p>
  </div>

  <div class="field">
    <label class="label">City</label>
    <div class="control">
      <input bind:value={profile.city} class="input" type="text" placeholder="Krasnodar">
    </div>
  </div>

  <div class="field">
  	<label class="label">Interests</label>
  	<div class="control">
  		<input class="input" type="text" data-type="tags" placeholder="Choose Tags">
  	</div>
  </div>

  <div class="field is-grouped">
    <div class="control">
      <button class="button is-link" type="button" on:click={submitProfile(profile)}>Submit</button>
    </div>
  </div>
</div>

<style>
.container {
  border: 1px solid white;
  border-radius: 14px;
  padding: 25px;
  background-color: #f1f4f7;
  margin-top: 12vh;
}
</style>
