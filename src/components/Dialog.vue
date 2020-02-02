<template>
  <v-dialog v-model="open" persistent max-width="600px">
    <v-card>
      <v-card-text>
        <v-container>
          <v-row>
            <v-col cols="12">
              <v-text-field
                label="Title *"
                v-model="title"
                :rules="[rules.required]"
                required
              ></v-text-field>
            </v-col>
            <v-col cols="12">
              <v-textarea v-model="text" color="teal">
                <template v-slot:label>
                  <div>
                    Text
                    <small>(optional)</small>
                  </div>
                </template>
              </v-textarea>
            </v-col>
            <v-col cols="12" sm="6">
              <v-select
                v-model="color"
                :items="colors"
                label="Color"
              ></v-select>
            </v-col>
          </v-row>
        </v-container>
        <small>*indicates required field</small>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="blue darken-1" text @click="$emit('close-dialog')">Close</v-btn>
        <v-btn
          color="blue darken-1"
          text
          :disabled="!title.trim()"
          @click="$emit('save-dialog', {color, text:text.trim(), title:title.trim()})"
        >Save</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  name: 'modal',
  props: {
    curTitle: String,
    curText: String,
    curColor: String,
    open: Boolean,
  },
  data() {
    const colorList = ['indigo', 'purple', 'deep-orange', 'teal', 'brown'];
    return {
      colors: colorList,
      color: this.curColor || colorList[0],
      text: this.curText,
      title: this.curTitle,
      rules: {
        required: value => !!value.trim() || 'Required',
      },
    };
  },
  watch: {
    curColor(newValue) {
      this.color = newValue;
    },
    curText(newValue) {
      this.text = newValue;
    },
    curTitle(newValue) {
      this.title = newValue;
    },
  },
};
</script>

<style scoped>
</style>
