# import re
# import string
# from torch import clamp
# from transformers import AutoTokenizer, AutoModel
# from sklearn.metrics.pairwise import cosine_similarity

# class TokenSimilarity:

#     def load_pretrained(self, from_pretrained:str="indobenchmark/indobert-base-p1"):
#         self.tokenizer = AutoTokenizer.from_pretrained(from_pretrained)
#         self.model = AutoModel.from_pretrained(from_pretrained)
        
#     def __cleaning(self, text:str):
#         # clear punctuations
#         text = text.translate(str.maketrans('', '', string.punctuation))

#         # clear multiple spaces
#         text = re.sub(r'/s+', ' ', text).strip()

#         return text
        
#     def __process(self, first_token:str, second_token:str):
#         inputs = self.tokenizer([first_token, second_token],
#                                 max_length=self.max_length,
#                                 truncation=self.truncation,
#                                 padding=self.padding,
#                                 return_tensors='pt')

#         attention = inputs.attention_mask

#         outputs = self.model(**inputs)

#         # get the weights from the last layer as embeddings
#         embeddings = outputs[0] # when used in older transformers version
#         # embeddings = outputs.last_hidden_state # when used in newer one

#         # add more dimension then expand tensor
#         # to match embeddings shape by duplicating its values by rows
#         mask = attention.unsqueeze(-1).expand(embeddings.shape).float()

#         masked_embeddings = embeddings * mask
        
#         # MEAN POOLING FOR 2ND DIMENSION
#         # first, get sums by 2nd dimension
#         # second, get counts of 2nd dimension
#         # third, calculate the mean, i.e. sums/counts
#         summed = masked_embeddings.sum(1)
#         counts = clamp(mask.sum(1), min=1e-9)
#         mean_pooled = summed/counts

#         # return mean pooling as numpy array
#         return mean_pooled.detach().numpy()
        
#     def predict(self, first_token:str, second_token:str,
#                 return_as_embeddings:bool=False, max_length:int=16,
#                 truncation:bool=True, padding:str="max_length"):
#         self.max_length = max_length
#         self.truncation = truncation
#         self.padding = padding

#         first_token = self.__cleaning(first_token)
#         second_token = self.__cleaning(second_token)

#         mean_pooled_arr = self.__process(first_token, second_token)
#         if return_as_embeddings:
#             return mean_pooled_arr

#         # calculate similarity
#         similarity = cosine_similarity([mean_pooled_arr[0]], [mean_pooled_arr[1]])

#         return similarity

# token1 = 'jalan berlubang akibat kendaraan bermuatan berat mohon segera ada perbaikan mengingat kalu malam penerangan kurang. Kelurahan Karangwaru, Kecamatan Tegalrejo, Daerah Istimewa Yogyakarta'
# token2 = 'Jalan umum Amblong RT 04 RW 02 J Gotong Royong utara Kelurahan Karangwaru Kemantren Tegalrejo Kota Yogyakarta. Kelurahan Karangwaru, Kecamatan Tegalrejo, Daerah Istimewa Yogyakarta'
# token3 = 'jalan umum wilayah rt 04 r 02, jalan ambles sangat membahayakan pengguna jalan. Kelurahan Karangwaru, Kecamatan Tegalrejo, Daerah Istimewa Yogyakarta'

# print(model.predict(token1, token2))
# print(model.predict(token1, token3))

import re
import string
from Sastrawi.StopWordRemover.StopWordRemoverFactory import StopWordRemoverFactory
from torch import clamp
from transformers import AutoTokenizer, AutoModel
from sklearn.metrics.pairwise import cosine_similarity

class TokenSimilarity:

    def load_pretrained(self, from_pretrained:str):
        self.tokenizer = AutoTokenizer.from_pretrained(from_pretrained)
        self.model = AutoModel.from_pretrained(from_pretrained)
        
    def cleaning(self, text:str):
        # casefolding
        text = text.lower()

        # stopwords
        srf = StopWordRemoverFactory()
        stop = srf.create_stop_word_remover()
        text = stop.remove(text)

        return text
        
    def process(self, first_token:str, second_token:str):
        inputs = self.tokenizer([first_token, second_token],
                                max_length=self.max_length,
                                truncation=self.truncation,
                                padding=self.padding,
                                return_tensors='pt')

        attention = inputs.attention_mask

        outputs = self.model(**inputs)

        # get the weights from the last layer as embeddings
        #embeddings = outputs[0] # when used in older transformers version
        embeddings = outputs.last_hidden_state # when used in newer one

        # add more dimension then expand tensor
        # to match embeddings shape by duplicating its values by rows
        mask = attention.unsqueeze(-1).expand(embeddings.shape).float()

        masked_embeddings = embeddings * mask
        
        # MEAN POOLING FOR 2ND DIMENSION
        # first, get sums by 2nd dimension
        # second, get counts of 2nd dimension
        # third, calculate the mean, i.e. sums/counts
        summed = masked_embeddings.sum(1)
        counts = clamp(mask.sum(1), min=1e-9)
        mean_pooled = summed/counts

        # return mean pooling as numpy array
        return mean_pooled.detach().numpy()
        
    def predict(self, first_token:str, second_token:str,
                return_as_embeddings:bool=False, max_length:int=16,
                truncation:bool=True, padding:str="max_length"):
        self.max_length = max_length
        self.truncation = truncation
        self.padding = padding

        first_token = self.cleaning(first_token)
        second_token = self.cleaning(second_token)

        mean_pooled_arr = self.process(first_token, second_token)
        if return_as_embeddings:
            return mean_pooled_arr

        # calculate similarity
        similarity = cosine_similarity([mean_pooled_arr[0]], [mean_pooled_arr[1]])

        return similarity