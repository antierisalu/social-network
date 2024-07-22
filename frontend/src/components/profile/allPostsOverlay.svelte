<script>
    import { createEventDispatcher } from "svelte";
    import { slide, fade } from "svelte/transition";
    import Posts from "../posts/posts.svelte"
    import Button from "../../shared/button.svelte";
  
    const dispatch = createEventDispatcher();
    function closeOverlay() {
      dispatch("close");
    }
    export let posts;
    console.log("HUUUUH", posts)
  
  </script>
  
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div out:fade={{ duration: 150 }} class="overlay"></div>
  <div in:slide out:fade={{ duration: 150 }} class="modal">
    <div class="modal-content">
        <Posts posts={posts} allowCreate={false}></Posts>
        <Button on:click={closeOverlay}>Close</Button>
    </div>
  </div>
  
  <style>
    .overlay {
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: min-content;
      background-color: rgba(0, 0, 0, 0.5);
      z-index: 0;
    }
  
    .modal {
      position: fixed;
      top: 50%;
      left: 50%;
      border: 1px solid black;
      transform: translate(-50%, -50%);
      background-color: #011;
      padding: 20px;
      box-shadow: 0 5px 15px rgba(0, 0, 0, 0.3);
      z-index: 20;
      width: 100%;
      max-width: 1000px;
      border-radius: 8px;
    }
  
    .modal-content {
      margin-top:20px;
      position: relative;
      overflow-y: scroll;
      max-height: 80vh;
      
    }
  </style>
  