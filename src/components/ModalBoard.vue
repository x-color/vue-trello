<template>
  <v-dialog v-model="open" persistent max-width="600px" @click:outside="close">
    <v-card>
      <v-container>
        <v-row>
          <!-- PIN: Title -->
          <v-col cols="10">
            <v-card-title class="py-0">
              <v-text-field
                class="title pt-1 pb-0"
                dense
                hide-details
                placeholder="Title..."
                v-model="value.title"
                autofocus
              ></v-text-field>
            </v-card-title>
          </v-col>

          <!-- PIN: Color -->
          <v-col cols="2" class="pl-0">
            <!-- PIN: Colors Menu -->
            <v-menu
              v-model="isOpenedMenu"
              offset-y
              :close-on-content-click="false"
            >
              <template v-slot:activator="{ on }">
                <v-icon
                  :color="iconColor(value.color)"
                  class="ma-1"
                  x-large
                  v-on="on"
                >{{ circleIcon(value.color) }}</v-icon>
              </template>
              <v-card>
                <v-list dense subheader max-width="400">
                  <v-subheader>Color</v-subheader>
                  <v-container>
                    <v-row>
                      <v-col
                        v-for="(color, index) in colors"
                        :key="index"
                        cols="auto"
                        class="pa-0"
                      >
                        <v-list-item class="px-0">
                          <v-icon
                            large
                            :color="iconColor(color)"
                            @click.stop="selectColor(color)"
                          >{{ circleIcon(color) }}</v-icon>
                        </v-list-item>
                      </v-col>
                    </v-row>
                  </v-container>
                </v-list>
              </v-card>
            </v-menu>
          </v-col>

          <!-- PIN: Text -->
          <v-col cols="12" class="pt-0">
            <v-card-text class="pt-0 pb-2">
              <v-textarea
                class="body-2"
                v-model="value.text"
                placeholder="Add desription..."
                auto-grow
                dense
                hide-details
              ></v-textarea>
            </v-card-text>
          </v-col>
        </v-row>
      </v-container>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn text color="red" @click="close">CANCEL</v-btn>
        <v-btn text color="green" :disabled="!value.title" @click="save">CONFIRM</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  name: 'ModalBoard',
  filters: {
    replaceToHintText(text) {
      if (!text) {
        return 'Add description...';
      }
      return text;
    },
    replaceToHintTitle(text) {
      if (!text) {
        return 'No title...';
      }
      return text;
    },
  },
  props: {
    value: Object,
    open: Boolean,
  },
  data() {
    return {
      editTitleMode: false,
      editTextMode: false,
      isOpenedMenu: false,
      colors: ['white', 'indigo', 'red', 'orange', 'pink', 'lime'],
    };
  },
  methods: {
    circleIcon(color) {
      if (color === 'white') {
        return 'mdi-circle-outline';
      }
      return 'mdi-circle';
    },
    iconColor(color) {
      if (color === 'white') {
        return 'gray';
      }
      return color;
    },
    selectColor(color) {
      this.value.color = color;
    },
    close() {
      if (!this.isOpenedMenu) {
        this.$emit('close');
      }
    },
    save() {
      if (!this.isOpenedMenu) {
        this.$emit('save');
      }
    },
  },
};
</script>
