from math import sin, cos, sqrt, atan2, radians

def get_distance_geolocations(lat1, lat2, lon1, lon2):
    R = 6373.0
    lat1 = radians(lat1)
    lon1 = radians(lat2)
    lat2 = radians(lon1)
    lon2 = radians(lon2)

    dlon = lon2 - lon1
    dlat = lat2 - lat1

    a = sin(dlat / 2)**2 + cos(lat1) * cos(lat2) * sin(dlon / 2)**2
    c = 2 * atan2(sqrt(a), sqrt(1 - a))

    distance = R * c /1000 

    return distance

def predict(first_text, second_text, model):
    # model predict
    # print("load model start")
    # model = TokenSimilarity()
    # model.load_pretrained('indobenchmark/indobert-base-p2')
    print("predict start")
    accuracy = float(model.predict(first_text, second_text))
    print("acc: ", accuracy)
    # return accuracy[0][0]
    return accuracy

def my_function(pred_text, source_text):
    acc = predict(pred_text, source_text)
    res = {
        'text': pred_text,
        'accuracy': acc
    }
    return res

# def predict_batch(text_batch, source_text):
#     # print("Parallel")
#     # num_cores = multiprocessing.cpu_count() / 2
#     num_cores = 2
#     processed_list = Parallel(n_jobs=num_cores)(delayed(my_function)(t, source_text) for t in text_batch)
#     return processed_list

# def main():
#     model = TokenSimilarity()
#     model.load_pretrained('indobenchmark/indobert-base-p2')
#     text_batch = ["halo halo", "halo 123", "halo broo", "halo halo", "halo 123", "halo broo", "halo halo", "halo 123", "halo broo"
#                 ,"halo halo", "halo 123", "halo broo", "halo halo", "halo 123", "halo broo", "halo halo", "halo 123", "halo broo"]
#     original_text = "wkowekowe"
#     start = datetime.now()
#     print("Parallel: ")
#     print(predict_batch(text_batch, original_text, model))
#     end = datetime.now()
#     print("Time taken: ", end-start)

    # print("Non Parallel: ")
    # start = datetime.now()
    # for text in text_batch:
    #     print(my_function(text, original_text, model))
    # end = datetime.now()
    # print("Time taken: ", end-start)

# main()