<script>
  import Button from "../../shared/button.svelte";
  import ImagePreview from "../../shared/imagePreview.svelte";


</script>

<main>
   
    <div class="post">Posts are coming</div>
    <div class="post">Posts are coming</div>
    <div class="post">Posts are coming</div>
    <div class="post">Posts are coming</div>
    <div class="post">Posts are coming
        <div class="comment">This is comment</div>
        <div class="comment">This is comment</div>
        <div class="comment">This is comment</div>
    </div>
    <div class="post">Posts are coming</div>
    <div class="post">Posts are coming</div>
    <div class="post">Posts are coming</div>
    <div class="post">Posts are coming</div>
    
    <div class="post">
            <div class="comment">This is comment</div>
            <div class="comment">This is comment</div>
            <div class="comment">This is comment</div>
    </div>
    <div class="createPost">
        <textarea placeholder="New post here.."></textarea>
        <div class="postButtons">
            <Button type="secondary">Post</Button>
            <Button type="secondary" inverse={true}>Add image</Button>
        </div>
    </div>
</main>

<style>

    main {
        padding: 4px;
        border-radius: 8px;
        border: solid 1px #333
    }

    div {
        border: solid 1px #333;
        padding: 8px;
        border-radius: 8px;
    }
    .postButtons {
        display: flex;
        flex-direction: column;
        padding: 0;
        margin-top:0.5em;
        border:none
    }

    .createPost {
        display: flex;
        align-items: center;
    }

    .createPost textarea {
        margin-right: 4px;
        flex-grow: 1;
        border-radius: 8px;
        min-height: 80px;
        margin-bottom: 0;
        padding: 4px;
        text-align: left;
        vertical-align: top;
        resize: none;
        overflow: auto;
    }

    .comment {
        margin: 2px;
    }
    .post {
        text-align: left;
        margin: 4px
    }
</style>