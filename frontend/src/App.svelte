<script>
  import { Router, Route } from "svelte-routing";
  import Home from "./routes/Home.svelte";
  import Docs from "./routes/Docs.svelte";
  import NotFound from "./routes/NotFound.svelte";
  import ViewDoc from "./lib/view_doc.svelte";
  import "@fontsource/inter";
  import Nav from "./routes/Nav.svelte";

  export let url = "";
</script>

<div class="container">
  <Router {url}>
    <Nav />

    <div class="content">
      <Route path="/"><Home /></Route>
      <Route path="/docs"><Docs /></Route>
      <Route path="/docs/:docID" component={ViewDoc} />
      <Route><NotFound /></Route>
    </div>
  </Router>
</div>

<style>
  .container {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .content {
    width: 100%;
    max-width: 800px;
    margin: 1rem auto;
  }
</style>
