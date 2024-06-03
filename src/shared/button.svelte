<script>
    export let type = 'primary'
    export let flat = false
    export let inverse = false
    export let customStyle = ''
    export let w84 = false

</script>

<button class={type} class:w84={w84} class:flat={flat} class:inverse={inverse} style={customStyle} on:click>
    <slot></slot>
</button>

<style>
    button {
        border: 0;
        cursor: pointer;
        border-radius: 8px;
        padding: 8px 12px;
        box-shadow: 1px 2px 3px rgba(0, 0, 0, 0.2);
    }
    .btn200px {
        width: 200px;
    }

    .w84 {
        width: 84px
    }

    .primary {
        background:crimson;
        color: darkslategrey;   
    }
    .secondary {
        background:greenyellow ;
        color: #011
    }
    .flat {
        box-shadow: none;
    }

    .primary.inverse {
        color: crimson;
        background: #011;
        border: 2px solid crimson;
    }

    .secondary.inverse:active{
        background-color: #999;
    }
    .primary.inverse:active{
        background-color: #999;
    }

    .secondary.inverse {
        color: greenyellow;
        background-color: #011;
        border: 2px solid rgb(105, 155, 30)
    }
    
</style>