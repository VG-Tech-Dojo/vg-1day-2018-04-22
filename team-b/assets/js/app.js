(function() {
  'use strict';
  const Message = function() {
    this.body = ''
    this.username = ''
  };

  Vue.component('message', {
    // 1-1. ユーザー名を表示しよう
    props: ['id', 'body', 'username', 'removeMessage', 'updateMessage'],
    data() {
      return {
        editing: false,
        editedBody: null,
      }
    },
    // 1-1. ユーザー名を表示しよう
    template: `
    <div class="message">
      <div v-if="editing">
        <div class="row">
          <textarea v-model="editedBody" class="u-full-width"></textarea>
          <button v-on:click="doneEdit">Save</button>
          <button v-on:click="cancelEdit">Cancel</button>
        </div>
      </div>
      <div class="message-body" v-else>
        <p>{{ username }}</p>
        <span>{{ body }}</span>
        <span class="action-button u-pull-right" v-on:click="edit">&#9998;</span>
        <span class="action-button u-pull-right" v-on:click="remove">&#10007;</span>
      </div>
    </div>
  `,
    methods: {
      remove() {
        this.removeMessage(this.id)
      },
      edit() {
        this.editing = true
        this.editedBody = this.body
      },
      cancelEdit() {
        this.editing = false
        this.editedBody = null
      },
      doneEdit() {
        this.updateMessage({id: this.id, body: this.editedBody})
          .then(response => {
            this.cancelEdit()
          })
      }
    }
  });

  const app = new Vue({
    el: '#app',
    data: {
      messages: [],
      newMessage: new Message(),
      bots: [
        {key: 'mb', options: [' status', ' reset', ' 1 1', ' 1 2', ' 1 3', ' 2 1', ' 2 2', ' 2 3', ' 3 1', ' 3 2', ' 3 3']},
        {key: 'talk', options: []},
        {key: 'gacha', options: []},
        {key: 'omikuji', options: []}
      ],
      index: 0,
      key: ''
    },
    created() {
      this.getMessages();
      setInterval(() => {
        axios.get('/api/messages')
        .then(res => {
          this.messages = res.data.result
          this.messages.reverse()
        })
      }, 1000)
    },
    methods: {
      getMessages() {
        fetch('/api/messages').then(response => response.json()).then(data => {
          this.messages = data.result;
          this.messages.reverse()
        });
      },
      keydown(e) {
        if (e.keyCode === 13 && (e.altKey || e.shiftKey || e.ctrlKey)) {
          this.sendMessage()
        }
        if (e.keyCode === 9) {
          e.preventDefault();
          if (this.index === 0) {
            this.key = this.newMessage.body;
            if (this.bots.filter(bot => bot.key.substr(0, this.key.length) === this.key).length > 0) {
              this.index = 1;
              const username = this.newMessage.username;
              this.newMessage = {
                body: this.bots.filter(bot => bot.key.substr(0, this.key.length) === this.key)[0].key,
                username
              }
            }
          } else {
          }
        } else {
          this.index = 0;
        }
      },
      sendMessage() {
        const message = this.newMessage;
        fetch('/api/messages', {
          method: 'POST',
          body: JSON.stringify(message)
        })
          .then(response => response.json())
          .then(response => {
            if (response.error) {
              alert(response.error.message);
              return;
            }
            this.messages.push(response.result);
            this.clearMessage();
          })
          .catch(error => {
            console.log(error);
          });
      },
      removeMessage(id) {
        return fetch(`/api/messages/${id}`, {
          method: 'DELETE'
        })
        .then(response => response.json())
        .then(response => {
          if (response.error) {
            alert(response.error.message);
            return;
          }
          this.messages = this.messages.filter(m => {
            return m.id !== id
          })
        })
      },
      updateMessage(updatedMessage) {
        return fetch(`/api/messages/${updatedMessage.id}`, {
          method: 'PUT',
          body: JSON.stringify(updatedMessage),
        })
        .then(response => response.json())
        .then(response => {
            if (response.error) {
              alert(response.error.message);
              return;
            }
            const index = this.messages.findIndex(m => {
              return m.id === updatedMessage.id
            })
            Vue.set(this.messages, index, updatedMessage)
        })
      },
      clearMessage() {
        this.newMessage.body = '';
      }
    },
    mounted() {
    }
  });
})();
