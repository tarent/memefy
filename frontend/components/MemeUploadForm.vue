<template>
    <div class="box">
        <div class="box-header with-border">
            <h3 class="box-title">MEME Upload</h3>
        </div>

        <form role="form" @submit.prevent="upload">
            <div class="box-body">
                <message-success v-if="uploadStatus === 'success'">MEME erfolgreich gespeichert</message-success>
                <message-error v-if="uploadStatus === 'fail'">MEME konnte nicht gespeichert werden!</message-error>

                <div class="form-group" :class="{'has-error': $v.meme.name.$error}">
                    <label for="meme-name">Name</label>
                    <input class="form-control" id="meme-name" placeholder="MEME Name" v-model="meme.name">
                </div>
                <div class="form-group" :class="{'has-error': $v.meme.pic.$error}">
                    <label for="meme-pic">MEME</label>
                    <input id="meme-pic" type="file" @change="setMemePicture">
                </div>
                <div class="form-group" :class="{'has-error': $v.meme.sound.$error}">
                    <label for="meme-sound">Sound</label>
                    <input id="meme-sound" type="file" @change="setMemeSound">
                </div>
            </div>
            <!-- /.box-body -->

            <div class="box-footer">
                <button type="submit" class="btn btn-primary">Speichern</button>
            </div>
        </form>
    </div>
</template>

<script>
    import { required, minValue } from 'vuelidate/lib/validators'
    import { mapActions } from 'vuex'

    export default {
        name: "meme-upload-form",
        data() {
            return {
                meme: {
                    pic: null,
                    sound: null,
                    name: '',
                    meta: {
                        tags: ['cool']
                    }
                },
                uploadStatus: null
            }
        },
        validations: {
            meme: {
                name: {
                    required,
                },
                pic: {
                    required,
                },
                sound: {
                    required,
                },
            }
        },
        methods: {
            ...mapActions({
                saveMeme: 'meme/saveMeme'
            }),
            setMemePicture(event) {
                this.meme.pic = event.target.files[0]
            },
            setMemeSound(event) {
                this.meme.sound = event.target.files[0]
            },
            upload() {
                this.$v.meme.$touch();

                if(!this.$v.meme.$error) {
                    const formData = new FormData();
                    formData.append("name", this.meme.name);
                    formData.append("pic", this.meme.pic);
                    formData.append("sound", this.meme.sound);
                    formData.append("meta", JSON.stringify(this.meme.meta));

                    this.$axios.post("/meme/", formData)
                        .then((result) => {
                            this.uploadStatus = 'success'

                            this.saveMeme({
                                name: this.meme.name,
                                pic: `/meme/${result.data.pic}`,
                                sound: `/meme/${result.data.sound}`,
                                meta: this.meme.meta,
                            })
                        }, (error) => {
                            this.uploadStatus = 'fail'

                            console.log(error);
                        });
                }
            }
        }
    }
</script>

<style scoped>
</style>