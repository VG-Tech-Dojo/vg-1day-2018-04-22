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
        <span>{{ body }} - {{ username }}</span>
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
      newMessage: new Message()
    },
    created() {
      this.getMessages();
    },
    methods: {
      getMessages() {
        fetch('/api/messages').then(response => response.json()).then(data => {
          this.messages = data.result;
        });
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
        this.newMessage = new Message();
      }
    }
  });
})();
