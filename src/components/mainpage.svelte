<script>
  import Chat from "./chat.svelte";
  import Footer from "./footer.svelte";
  import Profile from "./profile.svelte";

</script>

<main>
<div id="leftSidebar">
    <Profile />
</div>
<div id="mainWindow">
KOIKMUU
</div>
<div id="rightSidebar">
    <Chat />
</div>
<div id="footer">
    <Footer />
</div>
    


</main>

<style>
main{
    display:grid;
    grid-template-columns: 300px auto 220px;
    grid-template-rows: auto 50px;
    height: 90vh;
}
div{
    border:solid;
    border-color: rgb(99, 158, 10);
}
#footer{
    grid-column: 2/4;
}
#leftSidebar{
    grid-row:1/3;
}


</style>