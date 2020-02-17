<template>
  <v-row justify="center">
    <v-dialog v-model="open" max-width="290" @click:outside="$emit('cancel')">
      <v-card>
        <v-card-title class="headline">{{ title }}</v-card-title>
        <v-card-text>{{ text }}</v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn color="red darken-1" text @click="$emit('cancel')">Cancel</v-btn>
          <v-btn color="green darken-1" text @click="$emit('confirm')">Confirm</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-row>
</template>
<script>
export default {
  name: 'ModalConfirm',
  props: {
    title: String,
    text: String,
    open: Boolean,
  },
  computed: {
    dialog: {
      get() {
        return this.open;
      },
      set() {},
    },
  },
};
</script>
