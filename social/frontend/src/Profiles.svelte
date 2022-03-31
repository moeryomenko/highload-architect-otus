<script>
	import { onMount } from "svelte";
	import { Card, CardTitle, CardSubtitle, CardText, Row, Col } from 'svelte-materialify';
	import InfiniteScroll from "./InfiniteScroll.svelte";

	let nextToken = '';
	let data = [];
	let newBatch = [];

	function getNext() {
		return nextToken === '' ? '': `?page_token=${nextToken}`;
	}

	async function fetchPage() {
		fetch(`/api/v1/profiles${getNext()}`, {
			method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${sessionStorage.getItem('access_token')}`,
			}}).then(async response => {
				const object = await response.json();
				newBatch = object.profiles;
				nextToken = object.next_token_page;
			});
	};

	onMount(() => { fetchPage(); })

	$: data = [
		...data,
		...newBatch,
	];
</script>


{#each data as profile}
<Card>
  <Row>
    <Col cols={8}>
      <CardTitle>{profile.first_name} {profile.last_name}</CardTitle>
      <CardSubtitle>{profile.age} years old</CardSubtitle>
      <CardText>From: {profile.city}<br>Interests: {profile.interests.join(', ')}</CardText>
    </Col>
  </Row>
</Card>
{/each}
<InfiniteScroll on:loadMore={() => {fetchPage()}} />
